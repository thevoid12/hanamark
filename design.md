# design doc of hanamark

### objective
- main objective is to take my bunch of markdown files and convert it into html files. add templated headers and footers
- static build or deploy it in a server which will watch for any new files, take the file if it is dropped in a folder and convert it into html and commit the code in github
- we will have a bunch of templates and fill the template and send back based on the type of configured file
### things we need to build
- [ ] configurable folder paths and file structure
- [ ] templates: header and footer for each type of page or all ie needed to be added in all pages
- [ ] two modes flag. static build or a server kind of mode where it keeps watching for files in the folder. the moment file drops in it does the markdown conversion process
- [ ] configurable adding css accordingly
