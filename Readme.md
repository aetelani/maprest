# Quick testing 
http -jv POST :8181/cities <<< '{"name":"Test","country":{"name":"fi"}}'
http -jv :8181/cities/Test
http -j :8181/cities countries==fi continents==foo
