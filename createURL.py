import random

outputPath = 'url.txt'
f = open(outputPath, 'w+')
prefixList = ['baidu', 'google', 'yahoo', 'microsoft', 'csdn', 'bilibili', 'jianshu',
              'youtube', 'reddit', 'facebook', 'bitcointalk', 'github', 'arxiv', 'arxiv-sanity',
              'overleaf', 'processon']
suffixList = [str(i) for i in range(1000)]

for _ in range(10):
	for prefix in prefixList:
	    for suffix in suffixList:
	        url = "https://www." + prefix + ".com/" + suffix + ".jpg\n"
	        num = random.randint(10, 10000)
	        for i in range(num):
	            f.write(url)

f.close()
