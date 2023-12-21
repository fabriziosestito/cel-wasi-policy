#!/usr/bin/env bats


@test "accept" {
  run kwctl run annotated-policy.wasm -r test_data/pod-creation-postgres.json -s test_data/settings.json

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]
}

@test "reject" {
  run kwctl run annotated-policy.wasm -r test_data/pod-creation-nginx.json -s test_data/settings.json

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
}
