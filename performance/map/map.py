import time

# print (time.time())
# exit()

dic = {}
n = 0
while(n < 10000):
    dic[n] = n
    n += 1
start = time.clock()
for v in dic:
    for vv in dic:
        if dic[v] == dic[vv]:
            break
print time.clock() - start

list = []
n = 0
while (n < 10000):
    list.append(n)
    n += 1
start = time.clock()
for v in list:
    for vv in list:
        if vv == v:
            break
print time.clock() - start







