Can be deployed as GAE module. Incoming requests are similar to requests to WIX statics.

It does HEAD request to Cloud Storage to find if we have this file, if not - pings Compute Engine to do thumb/resize and gives it back.

**Usage**:

* /media/zzz.jpg
 will render image from GCS via google frontend servers

* /exif/media/zzz.jpg will render JSON with EXIF.
