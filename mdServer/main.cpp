#include <iostream>
#include <thread>
#include "marketspi.h"

using namespace std;

int main(int argc, char *argv[])
{
    //for (int i=1;i<argc;i++){
    //    cout <<argv[i]<< endl;
    //}
    if (7>argc) return 0;
    MarketSpi spi(argv[1],argv[2],argv[3],argv[4],argv[5],argv[6]);
    //cout << "Hello World!" << endl;
    thread th1(&MarketSpi::receive,spi);
    th1.join();
    return 0;
}
