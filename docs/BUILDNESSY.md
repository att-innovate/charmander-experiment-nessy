Build analyzers, vector and Nessy 
---------------------------------


#### Build Analytics
Change your working directory back to the root of Charmander and start the build process

    cd ..
    ./bin/build_analytics


#### Build Analytics
Next, build Vector, pcp, and data collector

    ./bin/build_vector

#### Build Nessy
Now let's build nessy.

    ./experiments/nessy/bin/build

This command creates and deploys Docker images for nessy.
This process will take some time the first time you run it.


Everything is set up and ready to run! it is now time to [Run a Script][3] or [Run Tasks Manually][4]

[3]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/RUNSCRIPT.md
[4]: https://github.com/att-innovate/charmander-experiment-nessy/blob/master/docs/RUNMANUALLY.md