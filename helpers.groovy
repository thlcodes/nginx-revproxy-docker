def getEnvsFromFile(envfile, cb) {
  def envs = readYaml file: envfile
  envs.each{k,v ->
    cb(k,v)
  }
}

def getAppNameFromManifest(manifestFile) {
  def manifest = readYaml file: manifestFile
  def name = manifest.applications[0].name
  return name
}

return this