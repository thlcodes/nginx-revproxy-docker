// +build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	errorslib "dev.beta.audi/gorepo/lib_go_common/errors"
	greeterpb "dev.beta.audi/gorepo/lib_proto_models/golib/greeter"
	"dev.beta.audi/gorepo/lib_proto_models/golib/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntegrationTests", func() {

	const (
		knownUsersKey   = "knownUsers"
		newFallbackName = "Mini-Me"
	)

	var (
		receivedMessage *greeterpb.Message
		err             error
	)

	AfterEach(func() {
		// Reset the expectations of the mock servers after each test
		storageMock.Mock.ExpectedCalls = []*mock.Call{}
	})

	Describe("Hello", func() {
		Context("when the name is inappropriate", func() {
			It("should return an error", func() {
				receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: "Thomas"})
				Expect(err).To(HaveOccurred())
				Expect(receivedMessage).To(BeNil())
			})

			It("should return an error", func() {
				receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: "idiot"})
				Expect(err).To(HaveOccurred())
				Expect(receivedMessage).To(BeNil())
			})

			Specify("an http request that should return an http error", func() {
				requestBody, err := json.Marshal(map[string]string{"name": cfg.Greeter.InappropriateNames[1]})
				Expect(err).NotTo(HaveOccurred())

				address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
				response, err := http.Post(address, "application/json", bytes.NewBuffer(requestBody))
				defer func() { _ = response.Close }()

				Expect(err).NotTo(HaveOccurred())

				reply, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

				var message *greeterpb.Message
				err = json.Unmarshal(reply, &message)
				Expect(err).NotTo(HaveOccurred())
				Expect(message.Message).To(Equal("could not say hello to 'idiot'"))
			})
		})

		Context("when the name to greet is empty", func() {
			var (
				serializationError error
				knownUsersList     []string
				listBytes          []byte
			)

			JustBeforeEach(func() {
				listBytes, serializationError = json.Marshal(knownUsersList)
				Expect(serializationError).NotTo(HaveOccurred())
				mockedGetResult := &storage.KeyValuePair{Key: knownUsersKey, Value: &storage.Document{Value: string(listBytes)}}
				storageMock.On("Get", mock.Anything, mock.Anything).Return(mockedGetResult, nil)
			})

			JustAfterEach(func() {
				storageMock.AssertExpectations(GinkgoT())
			})

			Context("when an empty name has not been used before", func() {
				It("should greet the default fallback name", func() {
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)

					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{})

					Expect(receivedMessage.Message).To(Equal("Hello " + cfg.Greeter.DefaultName))
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("an http request that should greet the default fallback name", func() {
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)

					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", nil)
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("Hello " + cfg.Greeter.DefaultName))
				})
			})

			Context("when an empty name has been used before", func() {
				BeforeEach(func() {
					knownUsersList = []string{cfg.Greeter.DefaultName}
				})

				It("should greet the default fallback name again", func() {
					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{})

					Expect(receivedMessage.Message).To(Equal("Hello " + cfg.Greeter.DefaultName + " again"))
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("an http request that should greet the default fallback name again", func() {
					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", nil)
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("Hello " + cfg.Greeter.DefaultName + " again"))
				})
			})

			Context("when an new fallback name has been set", func() {
				BeforeEach(func() {
					_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: newFallbackName})
					Expect(err).NotTo(HaveOccurred())
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
				})

				AfterEach(func() {
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
					_, err = client.Reset(context.Background(), &empty.Empty{})
					Expect(err).NotTo(HaveOccurred())
				})

				Context("when the new fallback name has not been greeted yet", func() {
					It("should greet the new fallback name", func() {
						receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{})

						Expect(receivedMessage.Message).To(Equal("Hello " + newFallbackName))
						Expect(err).NotTo(HaveOccurred())
					})

					Specify("an http request that should greet the new fallback name", func() {
						address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
						response, err := http.Post(address, "application/json", nil)
						defer func() { _ = response.Close }()

						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(http.StatusOK))

						reply, err := ioutil.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						var message *greeterpb.Message
						err = json.Unmarshal(reply, &message)
						Expect(err).NotTo(HaveOccurred())
						Expect(message.Message).To(Equal("Hello " + newFallbackName))
					})
				})

				Context("when the new fallback name has been greeted already", func() {
					BeforeEach(func() {
						knownUsersList = []string{newFallbackName}
					})

					It("should greet the new fallback name again", func() {
						receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{})

						Expect(receivedMessage.Message).To(Equal("Hello " + newFallbackName + " again"))
						Expect(err).NotTo(HaveOccurred())
					})

					Specify("an http request that should greet the new fallback name again", func() {
						address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
						response, err := http.Post(address, "application/json", nil)
						defer func() { _ = response.Close }()

						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(http.StatusOK))

						reply, err := ioutil.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						var message *greeterpb.Message
						err = json.Unmarshal(reply, &message)
						Expect(err).NotTo(HaveOccurred())
						Expect(message.Message).To(Equal("Hello " + newFallbackName + " again"))
					})
				})
			})
		})

		Context("when the name to greet has been provided", func() {
			const nameToGreet = "Dr. Evil"

			var (
				serializationError error
				knownUsersList     []string
				listBytes          []byte
			)

			JustBeforeEach(func() {
				listBytes, serializationError = json.Marshal(knownUsersList)
				Expect(serializationError).NotTo(HaveOccurred())
				mockedGetResult := &storage.KeyValuePair{Key: knownUsersKey, Value: &storage.Document{Value: string(listBytes)}}
				storageMock.On("Get", mock.Anything, mock.Anything).Return(mockedGetResult, nil)
			})

			JustAfterEach(func() {
				storageMock.AssertExpectations(GinkgoT())
			})

			Context("when the maximum number of known users has been reached", func() {
				BeforeEach(func() {
					for counter := int64(0); counter < cfg.Greeter.MaxKnownUsers; counter++ {
						knownUsersList = append(knownUsersList, strconv.Itoa(int(counter)))
					}
				})

				AfterEach(func() {
					knownUsersList = []string{}
				})

				It("should return an error", func() {
					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: nameToGreet})

					Expect(err).To(HaveOccurred())
					Expect(receivedMessage).To(BeNil())
				})

				Specify("an http request that should return an http error", func() {
					requestBody, err := json.Marshal(map[string]string{"name": nameToGreet})
					Expect(err).NotTo(HaveOccurred())

					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", bytes.NewBuffer(requestBody))
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("could not add user"))
				})
			})

			Context("when updating the known users list fails", func() {
				BeforeEach(func() {
					mockedError := errorslib.NewClientError("could not update known users list", errors.New(""))
					storageMock.On("Set", mock.Anything, mock.Anything).Return(nil, mockedError)
				})

				It("should return an error", func() {
					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: nameToGreet})

					Expect(err).To(HaveOccurred())
					Expect(receivedMessage).To(BeNil())
				})

				Specify("an http request that should return an http error", func() {
					requestBody, err := json.Marshal(map[string]string{"name": nameToGreet})
					Expect(err).NotTo(HaveOccurred())

					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", bytes.NewBuffer(requestBody))
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("could not update known users list"))
				})
			})

			Context("when the name to greet has not been greeted yet", func() {
				It("should greet the name", func() {
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)

					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: nameToGreet})

					Expect(receivedMessage.Message).To(Equal("Hello " + nameToGreet))
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("an http request that should greet the name", func() {
					storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)

					requestBody, err := json.Marshal(map[string]string{"name": nameToGreet})
					Expect(err).NotTo(HaveOccurred())

					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", bytes.NewBuffer(requestBody))
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("Hello " + nameToGreet))
				})
			})

			Context("when the name to greet has been greeted already", func() {
				BeforeEach(func() {
					knownUsersList = []string{nameToGreet}
				})

				It("should greet the name again", func() {
					receivedMessage, err = client.Hello(context.Background(), &greeterpb.HelloRequest{Name: nameToGreet})

					Expect(receivedMessage.Message).To(Equal("Hello " + nameToGreet + " again"))
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("an http request that should greet the name again", func() {
					requestBody, err := json.Marshal(map[string]string{"name": nameToGreet})
					Expect(err).NotTo(HaveOccurred())

					address := "http://" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) + "/hello"
					response, err := http.Post(address, "application/json", bytes.NewBuffer(requestBody))
					defer func() { _ = response.Close }()

					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					reply, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					var message *greeterpb.Message
					err = json.Unmarshal(reply, &message)
					Expect(err).NotTo(HaveOccurred())
					Expect(message.Message).To(Equal("Hello " + nameToGreet + " again"))
				})
			})
		})
	})

	Describe("GetFallbackName", func() {
		Context("when the fallback name has not been set", func() {
			It("should return the default fallback name", func() {
				var fallbackName *greeterpb.Name
				fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
				Expect(fallbackName.Name).To(Equal(cfg.Greeter.DefaultName))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a new fallback name has been set", func() {
			It("should return the new fallback name", func() {
				_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: newFallbackName})
				Expect(err).NotTo(HaveOccurred())

				var fallbackName *greeterpb.Name
				fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
				Expect(fallbackName.Name).To(Equal(newFallbackName))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the new fallback name has been reset", func() {
			BeforeEach(func() {
				_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: newFallbackName})
				Expect(err).NotTo(HaveOccurred())
			})

			JustBeforeEach(func() {
				storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
				_, err = client.Reset(context.Background(), &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				storageMock.AssertExpectations(GinkgoT())
			})

			It("should return the default fallback name", func() {
				var fallbackName *greeterpb.Name
				fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
				Expect(fallbackName.Name).To(Equal(cfg.Greeter.DefaultName))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("SetFallbackName", func() {
		Context("when the fallback name is empty", func() {
			It("should return an error", func() {
				_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: ""})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the fallback name has the required length", func() {
			It("should set the name as new default name", func() {
				_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: newFallbackName})
				Expect(err).NotTo(HaveOccurred())

				var fallbackName *greeterpb.Name
				fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
				Expect(fallbackName.Name).To(Equal(newFallbackName))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Reset the fallback name", func() {
		AfterEach(func() {
			var fallbackName *greeterpb.Name
			fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
			Expect(fallbackName.Name).To(Equal(cfg.Greeter.DefaultName))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when no fallback name has been set", func() {
			It("should succeed", func() {
				storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
				_, err = client.Reset(context.Background(), &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				storageMock.AssertExpectations(GinkgoT())
			})
		})

		Context("when a new fallback name has been set", func() {
			BeforeEach(func() {
				_, err = client.SetFallbackName(context.Background(), &greeterpb.SetFallbackNameRequest{Name: newFallbackName})
				Expect(err).NotTo(HaveOccurred())
			})

			JustBeforeEach(func() {
				var fallbackName *greeterpb.Name
				fallbackName, err = client.GetFallbackName(context.Background(), &empty.Empty{})
				Expect(fallbackName.Name).To(Equal(newFallbackName))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should reset the fallback name", func() {
				storageMock.On("Set", mock.Anything, mock.Anything).Return(&empty.Empty{}, nil)
				_, err = client.Reset(context.Background(), &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				storageMock.AssertExpectations(GinkgoT())
			})
		})
	})
})
