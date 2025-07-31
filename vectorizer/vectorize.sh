#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
set -e
DOCUMENTS_DIR="${SCRIPT_DIR}/../downloader/documents"
INSTITUTIONS=$(ls "${DOCUMENTS_DIR}")


for INSTITUTION in ${INSTITUTIONS}; do
  INSTITUTION_DIR="${DOCUMENTS_DIR}/${INSTITUTION}"
  INSTITUTION_CODE=$(echo "${INSTITUTION}" | tr '[:lower:]' '[:upper:]')  
  for FILE in "${INSTITUTION_DIR}"/*.pdf; do
    if [ -f "${FILE}" ]; then
      BASENAME=$(basename "${FILE}")
      echo "Processing file: ${BASENAME} for institution: ${INSTITUTION_CODE}"
      ${SCRIPT_DIR}/bin/vectorizer "${FILE}" "${INSTITUTION_CODE}"
      STATUS=$?
      if [ $STATUS -ne 0 ]; then
        echo "Error processing file: ${BASENAME} for institution: ${INSTITUTION}"
      fi
    fi
  done
  echo "Finished processing files for ${INSTITUTION_CODE}."
done

