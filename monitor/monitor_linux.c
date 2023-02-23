#include <stdio.h>
#include <stdlib.h>

// TODO: use getloadavg() instead of /proc/loadavg

static int get_loadavg(double *avg) {
    FILE *fp;
    char buf[256];
    int n;

    if ((fp = fopen("/proc/loadavg", "r")) == NULL) {
        return -1;
    }

    if (fgets(buf, sizeof(buf), fp) == NULL) {
        fclose(fp);
        return -1;
    }

    fclose(fp);

    n = sscanf(buf, "%lf %lf %lf", avg, avg + 1, avg + 2);
    if (n != 3) {
        return -1;
    }

    return 0;
}

int main(int argc, char *argv[]) {
    double avg[3];

    if (get_loadavg(avg) < 0) {
        fprintf(stderr, "get_loadavg() failed");
        return -1;
    }

    printf("loadavg: %.2f %.2f %.2f\n", avg[0], avg[1], avg[2]);
    return 0;
}
