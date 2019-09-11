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
    ctpspi ctp(argv[1],argv[2],argv[3]);
    thread th3(&ctpspi::runMarket,&ctp,argv[5]);
    thread th4(&ctpspi::runTrader,&ctp,argv[4]);
    thread th1(&ctpspi::runMRecv,&ctp);
    thread th2(&ctpspi::runTRecv,&ctp);
    //cout << "4" << endl;
    th1.join();
    th2.join();
    th3.join();
    th4.join();
    cout << "Hello World!" << endl;
    return 0;
}
