package program

// Version can be requested through the command line with:
//
//     godiffsub -v
//
// How to release new versions:
//   1. On the master branch, commit the new version in main.go.
//      Normally we would never commit to master but this is a change
//      that should never break the build. Tip: You can edit the file directly here.
//   2. Tag the new HEAD with the version (all tags have a v prefix). Push the new tag to Github.
//   3. Once the tag is pushed it will show as an unpublished versions on the Github releases page.
//      You will now need to edit it, the release name will be the same as the tag.
//      The description should contain bullet points of each of the pull requests and issues resolves in this patch.
//
const Version = "v0.0.1 2018-05-08"

