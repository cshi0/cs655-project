import random
import hashlib
import requests

chars = [chr(i) for i in range(65, 91)] + [chr(i) for i in range(97, 123)]

strs = []

for _ in range(10):
  s = ""
  for _ in range(5):
    s += random.choice(chars)
  strs += [s]

print(strs)

for s in strs:
  h = hashlib.md5(s.encode())
  print("curl -X POST http://localhost:8080/crackTask --header \"Content-Type: application/json\" --data '{\"ToUnhash\": \"" +h.hexdigest()+ "\"}'")
  requests.post('http://localhost:8080/crackTask', {"ToUnhash": h.hexdigest()})