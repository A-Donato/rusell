#!/bin/bash
echo "Welcome to our tiny script - lets deploy our function"

gcloud functions deploy russell-scrap \
--gen2 \
--region=us-central1 \
--runtime=go122 \
--source=C:/Users/alexi/Documents/github/rusell-cloud-funtions/src \
--entry-point=start-scrapping \
--trigger-http