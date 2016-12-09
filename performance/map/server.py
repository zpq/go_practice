def g(n):
    while n > 0:
        yield n
        n -= 1

a = g(10)

# for i in range(4):
print (a.__next__())
print (a.send(None))
print (a.__next__())
print (a.send(3))
print (a.__next__())

exit()


def countdown(n):
    # print("Counting down from ", n)
    while n >= 0:
        print('se')
        newvalue = yield n
        if newvalue is not None:
            n = newvalue
        else:
            n -= 1

if __name__ == "__main__":
    c = countdown(5)
    for x in c:
        print(x)
        if x == 5:
            c.send(3)

def generateList1(start,stop):
	for i in range(start,stop):
		yield i

a = generateList1(0, 5)

for i in range(0, 5):
    print(a.__next__())
    
    
