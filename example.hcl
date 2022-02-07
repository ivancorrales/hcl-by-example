TEST_FOLDER = "__test__"
EXECUTION_ID = random(6)

job check "this is a temporal job" {
  python "run.py" {}
}

job e2e "running integration tests" {

  python "app-e2e.py" {
    root_dir = TEST_FOLDER
  }

  slack {
    channel  = "slack-my-channel"
    message = "Job execution ${EXECUTION_ID} completed successfully"
  }
}



