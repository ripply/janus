#!/bin/sh
doGithubWorkflowProcessing () {
  if [ "" != "$GITHUB_ACTION" ] ; then
    # running in a github action, output results for next workflow action
    PASSING=`docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) passing.*/::set-output name=PASSING=::\1/p'`
    FAILING=`docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) failing.*/::set-output name=FAILING=::\1/p'`
    PENDING=`docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) pending.*/::set-output name=PENDING=::\1/p'`

    if [ "" != "$PASSING" ] ; then
      # truffle will exit with a non-zero exit code if any tests fail
      # when running from a github workflow, we need to determine what constitutes "failure" ourselves here
      # "failure" being the workflow fails and an email is triggered about the failing job
      # right now there are many failing tests (100<x<200)

    fi
    docker cp ci_openzeppelin_1:/openzeppelin-contracts/output.json ./openzeppelin-contracts-output.json
    if [ -e ./openzeppelin-contracts-output.json ] ; then
      if [ -e ./openzeppelin ] ; then
            ./openzeppelin \
        --expected ./openzeppelin-contracts/expected.json \
        --input ./openzeppelin-contracts/output.json \
        --output ./openzeppelin-contracts/updated.json
      fi
    else
      echo "Can't find truffle result json"
      return -1
    fi
    
    return $?
  fi
}
cleanupDocker () {
  docker-compose -f docker-compose-openzeppelin.yml -p ci kill
  docker-compose -f docker-compose-openzeppelin.yml -p ci rm -f
}
trap 'cleanupDocker ; echo "Tests Failed For Unexpected Reasons"' HUP INT QUIT PIPE TERM
docker-compose -p ci -f docker-compose-openzeppelin.yml build && docker-compose -p ci -f docker-compose-openzeppelin.yml up -d
if [ $? -ne 0 ] ; then
  echo "Docker Compose Failed"
  exit 1
fi
docker logs ci_openzeppelin_1 -f&
EXIT_CODE=`docker wait ci_openzeppelin_1`
if [ -z ${EXIT_CODE+z} ] || [ -z ${EXIT_CODE} ] || ([ "0" != "$EXIT_CODE" ] && [ "" != "$EXIT_CODE" ]) ; then
  # docker logs qtum_seeded_testchain
  # docker logs ci_janus_1
  # docker logs ci_openzeppelin_1
  doGithubWorkflowProcessing
  echo "Tests Failed - Exit Code: $EXIT_CODE (truffle exit code indicates how many tests failed)"
else
  EXIT_CODE=doGithubWorkflowProcessing
  echo "Tests Passed"
fi
cleanupDocker
exit $EXIT_CODE
