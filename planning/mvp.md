# MVP

Goal 1: Service can collect, store and display most recent run
* Seperate and log every component that is part ot the prompt
* Design API
* Build base version that dumps donta to a dater volume
* Displays a basic UI of the data

Goal 2: Benchmark Web & Image Sumparisation dump acta san benchmark file
* Update current bench system to accept web
* Modify metrics
* Decide on new image benchmark model
* Build of update benchmarking system
* Create new script that parses benchmark data and downloads the feed data
* Remove inline dump support

Goal 3: Auto Benchmarking Support
* Build data storage system, including benchmark data + bench
* Allow for multiple runs to be stored
* Run benchmark on cron for all components, or on trigger
* Provide output in Web UI

Goal 4: Optimisation Mode
* Allow service to modify Persona + Prompt
* Use existing benchmark data o alla for outimisation of a single entry, then validation across many
* Add radamentry quality tracking over time chart

Goal 5: Deep Metries Visualisation
* Track prompt Version, model,Persona over time
* Catagorise entries based on quality and data sources
