#include"v8_export.h"
int main(){
    void* p=NewIsolate(1);
    IsolateDispose(p);
}