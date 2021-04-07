data_dir = "/consul/data"
log_level = "DEBUG"
server = true

bootstrap_expect = 1
ui = true

bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"

connect {
  enabled = true
}

enable_script_checks = true
enable_local_script_checks = true
enable_central_service_config = true