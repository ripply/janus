#!/bin/bash
EXPECTED_OUTPUT=openzeppelin-contracts-expected-output.json
RESULUT_OUTPUT=openzeppelin-contracts-result-output.json
PRUNED_OUTPUT=openzeppelin-contracts-pruned-output.json

doGithubWorkflowProcessing() {
  if [ "" != "$GITHUB_ACTION" ] ; then
    echo "Running within github actions... processing output file"
    # running in a github action, output results for next workflow action
    echo Parsing passing...
    docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) passing.*/::set-output name=PASSING=::\1/p'
    echo Parsing pending...
    docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) pending.*/::set-output name=PENDING=::\1/p'
    echo Parsing failing...
    docker logs ci_openzeppelin_1 | sed -n 's/.* \([0-9]\{1,\}\) failing.*/::set-output name=FAILING=::\1/p'

    if [ ! -f $EXPECTED_OUTPUT ] ; then
      echo "Expected output not found -" $EXPECTED_OUTPUT
      doGithubWorkflowProcessingResult=1
      return
    fi

    if [ -e $RESULUT_OUTPUT ] ; then
      echo Successfully copied output results from docker container
      docker run --rm -v `pwd`:/output qtum/janus-openzeppelin-test \
        --expected /output/$EXPECTED_OUTPUT \
        --input /output/$RESULUT_OUTPUT \
        --output /output/$PRUNED_OUTPUT
    else
      echo "Failed to find output results in docker container"
      doGithubWorkflowProcessingResult=-1
      return
    fi
    
    doGithubWorkflowProcessingResult=$?
  else
    echo "Not running within github actions"
  fi
}
cleanupDocker () {
  echo Shutting down docker-compose containers
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
echo "Processing openzeppelin test results with exit code of:" $EXIT_CODE
doGithubWorkflowProcessingResult=$EXIT_CODE

if [ -e $RESULUT_OUTPUT ] ; then
  echo "Deleting existing output results"
  rm $RESULUT_OUTPUT
fi

echo "Copying output results from docker container to local filesystem"
docker cp ci_openzeppelin_1:/openzeppelin-contracts/output.json $RESULUT_OUTPUT

doGithubWorkflowProcessing
EXIT_CODE=$doGithubWorkflowProcessingResult
if [ -z ${EXIT_CODE+z} ] || [ -z ${EXIT_CODE} ] || ([ "0" != "$EXIT_CODE" ] && [ "" != "$EXIT_CODE" ]) ; then
  # these logs are so large we can't print them out into github actions
  # docker logs qtum_seeded_testchain
  # docker logs ci_janus_1
  # docker logs ci_openzeppelin_1
  echo "Tests Failed - Exit Code: $EXIT_CODE (truffle exit code indicates how many tests failed)"
else
  echo "Tests Passed"
fi
cleanupDocker
exit $EXIT_CODE
