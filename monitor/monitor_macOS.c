#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/sysctl.h>
#include <sys/resource.h>

// sysctl -n vm.loadavg

int main(int argc, char *argv[])
{
    double avg[3];

    if (getloadavg(avg, 3) < 0)
    {
        perror("getloadavg");
        return -1;
    }
    // 1, 5 and 15 minute load averages
    printf("loadavg: %.2f %.2f %.2f\n", avg[0], avg[1], avg[2]);

    struct rusage usage;
    while (1)
    {
        int ret = getrusage(RUSAGE_SELF, &usage);
        if (ret < 0)
        {
            perror("getrusage");
            exit(1);
        }

        printf("ru_maxrss: %ld\n", usage.ru_maxrss);
        printf("ru_ixrss: %ld\n", usage.ru_ixrss);
        printf("ru_idrss: %ld\n", usage.ru_idrss);
        printf("ru_isrss: %ld\n", usage.ru_isrss);
        printf("ru_minflt: %ld\n", usage.ru_minflt);
        printf("ru_majflt: %ld\n", usage.ru_majflt);
        printf("ru_nswap: %ld\n", usage.ru_nswap);
        printf("ru_inblock: %ld\n", usage.ru_inblock);
        printf("ru_oublock: %ld\n", usage.ru_oublock);
        printf("ru_msgsnd: %ld\n", usage.ru_msgsnd);
        printf("ru_msgrcv: %ld\n", usage.ru_msgrcv);
        printf("ru_nsignals: %ld\n", usage.ru_nsignals);
        printf("ru_nvcsw: %ld\n", usage.ru_nvcsw);
        printf("ru_nivcsw: %ld\n", usage.ru_nivcsw);

        printf("ru_utime: %ld.%06d\n", usage.ru_utime.tv_sec, usage.ru_utime.tv_usec);
    }

    return 0;
}