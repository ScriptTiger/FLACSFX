[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# FLACSFX: FLAC SelF-eXtracting archive
FLACSFX is a minimal FLAC-to-WAV transcoder and multitrack FLAC mixer for embedded FLAC archives. The FLAC files are embedded and archived by appending them to a precompiled FLACSFX executable. The transcoder component transcodes embedded FLAC files to WAV files. The multitrack FLAC mixer can mix any combination of compatible embedded FLAC tracks and output the mix either to a file or standard output, to be played by VLC or other compatible applications without needing to write it to a file.

Usage: `flacsfx [options...]`
Argument                  | Description
--------------------------|-----------------------------------------------------------------------------------------------------
 `-i <#,#-#,...>`         | Index entries to include
 `-o <directory\|file>`   | Destination directory, or file for single entry or mix
 `-flac`                  | Output FLAC
 `-mix`                   | Output mix
 `-b <16\|24\|32>`        | Bit depth of mix
 `-info`                  | Show stream info

By default, a new directory is created to extract the audio files to. However, if only one entry is selected, or a mix is being created, the output is a file and not a directory. `-` can thus be used in place of `<file>` to designate standard output as the destination for a single entry or mix, allowing single entries and mixes to be sampled by compatible software, such as VLC, without needing to write a file.

Without any arguments, the embedded FLACs will be transcoded into the working directory within a new directory of the same name as the executable, or into the working directory directly if only one embedded FLAC is present. The extracted audio files will either be titled using the pattern of `#_<executable name>`, or titled after the TITLE metadata tag if present within the FLAC. So, command-line usage is only optional and the end user can just execute the application as they would any other application for this default behavior.

# Appending FLAC files to a FLACSFX executable
Download the latest precompiled releases for the intended target system:  
https://github.com/ScriptTiger/FLACSFX/releases/latest

For appending FLAC files to a FLACSFX executable, issue one of the following commands.

For Windows:
```
copy /b "FLACSFX.exe"+"file1.flac"+"file2.flac"+... "MyFLACSFX.exe"
```

For Linux and Mac:
```
cat "FLACSFX" "file1.flac" "file2.flac" ... > "MyFLACSFX"
```

# Single-track FLACSFX
If you don't need any of the more advanced features of the newer version, such as built-in mixing and multitrack capabilities for embedding multiple audio files, the older version for simply embedding a single audio file is still available, although no longer maintained.

README:  
https://github.com/ScriptTiger/FLACSFX/tree/4ffca93e97a915359c7511f26735ca35a7eac0bd

Precompiled executables:  
https://github.com/ScriptTiger/FLACSFX/releases/tag/25AUG2024

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=MZ4FH4G5XHGZ4)
