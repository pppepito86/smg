#!/bin/bash

counter=1
for f in *.in; do 
	echo "$f"
	mv  "$f" "input$counter"
	counter=$((counter+1))
done

counter=1
for f in *.sol; do 
	echo "$f"
	mv  "$f" "output$counter"
	counter=$((counter+1))
done



