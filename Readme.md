# Build
docker build . && docker run <image>

# Quick testing 
http -jv POST :8181/cities <<< '{"name":"Test","country":{"name":"fi"}}'

http -jv :8181/cities/Test

http -j :8181/search/cities countries==fi continents==foo

...

see Dockerfile for building

Testing with client: go test -v

Todo Add swagger doc
