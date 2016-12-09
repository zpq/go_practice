import time

n = 0
list = []

while (n < 100000):
    list.append(n)
    n += 1
print len(list)

start = time.clock()
for v in list:
    if v in list:
        pass
print time.clock() - start