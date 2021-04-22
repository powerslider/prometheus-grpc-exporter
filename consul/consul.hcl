//data_dir = "/consul/data"
//log_level = "DEBUG"
server = true
////
////bootstrap_expect = 1
//bind_addr = "0.0.0.0"
//client_addr = "0.0.0.0"
enable_local_script_checks = true
//enable_central_service_config = true
//
ui_config {
  enabled = true
}

ports {
  grpc = 8502
}

connect {
  enabled = true
}