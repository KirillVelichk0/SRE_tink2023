#!/bin/bash
cd ../.. | sudo docker stack services prod | grep prod_oncall-green