module.exports = function (grunt) {
  grunt.loadNpmTasks("grunt-release");

  grunt.initConfig({
    release: {
      options: {
        npm: false,
        file: 'bower.json'
      }
    }
  });

};

