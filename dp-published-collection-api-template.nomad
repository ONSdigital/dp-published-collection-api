job "dp-published-collection-api" {
  datacenters = ["DATA_CENTER"]
  update {
    stagger = "10s"
    max_parallel = 1
  }
  group "dp-published-collection-api" {
    task "dp-published-collection-api" {
      artifact {
        source = "s3::S3_TAR_FILE"
        destination = "."
        // The Following options are needed if no IAM roles are provided
        // options {
        // aws_access_key_id = ""
        // aws_access_key_secret = ""
        // }
      }
      env {
        PORT = "${NOMAD_PORT_http}"
        DB_ACCESS = "PUBLISH_DATABASE_URL"
      }
      driver = "exec"
      config {
        command = "bin/dp-published-collection-api"
      }
      resources {
        cpu = 500
        memory = 350
        network {
          port "http" {}
        }
      }
      
    }
  }
}