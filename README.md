# dockerif
docker for inspectfile

Contains Dockerfile and all the code for inspectfile, a swiss tool to inspect a file with the following applications :
  - Siegfried
  - Fido
  - Mediainfo
  - Exiftool
  
To shrink the size of the container, the ClamAV is removed and the apt-get is optimisez with the options --no-install-recommends. 

inspectfile can be used with API calls to inspect local files, with curl to inspect a file POST with the call.

An interactive web page is also implemented to test inspectfile directly in your browser which inspect the file by drag and drop.



