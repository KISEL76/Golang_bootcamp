#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *ask_cow(char phrase[]) {
    int phrase_len = strlen(phrase);
    char *buf = (char *)malloc(sizeof(char) * (160 + (phrase_len + 2) * 3));
    strcpy(buf, " ");

    for (int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "-");
    }
    strcat(buf, "\n< ");
    strcat(buf, phrase);
    strcat(buf, " >\n ");
    for (int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "-");
    }
    strcat(buf, "\n        \\   ^__^\n");
    strcat(buf, "         \\  (oo)\\_______\n");
    strcat(buf, "            (__)\\       )\\/\\\n");
    strcat(buf, "                ||----w |\n");
    strcat(buf, "                ||     ||\n");

    return buf;
}