Can be deployed as GAE module. Incoming requests are similar to requests to WIX statics.

Goal is to ping cloud storage with HEAD to find if we have this file, if not - ping compute engine to do thumb/resize and give it back.

**Usage**:

* /media/zzz.jpg
 will render image itself from GCS via google frontend servers

* /exif/media/zzz.jpg will render JSON with EXIF.
