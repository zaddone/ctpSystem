mkdir traderServer/build
cd traderServer/build
cmake ../
make
cp traderServer /data/ctp/bin/
cd ../../
mkdir mdServer/build
cd mdServer/build
cmake ../
make
cp mdServer /data/ctp/bin/
cd ../../

