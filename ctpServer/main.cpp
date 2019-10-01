#include <iostream>
#include <thread>
#include "ctpspi.h"
using namespace std;

int main(int argc, char *argv[])
{

    //cout << argc << endl;
    //cout << argv[0] << endl;
    //for (int i=1;i<argc;i++){
    //    cout <<argv[i]<< endl;
    //}
    ctpspi *ctp = new ctpspi(argv[1],argv[2],argv[3]);
    //ctpspi *ctp = new ctpspi();
    thread th3(&ctpspi::runMarket,ctp,argv[5]);
    thread th4(&ctpspi::runTrader,ctp,argv[4]);
    thread th1(&ctpspi::runMRecv,ctp);
    thread th2(&ctpspi::runTRecv,ctp);
    //ctp->getConfigM();
    //ctp->getConfigT();
    //cout << "4" << endl;
    th1.join();
    th2.join();
    th3.join();
    th4.join();
    cout << "Hello World!" << endl;
    delete ctp;
    return 0;
}
