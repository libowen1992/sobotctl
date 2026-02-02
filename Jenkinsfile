pipeline{
    agent any
    options{
	    disableConcurrentBuilds()
	}
    environment {
        COS_URl="https://v6-customer-1304925624.cos.ap-beijing.myqcloud.com/deploy/"
        PATH="/data/service/node/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/data/service/jdk/bin:/data/service/maven/bin:/data/service/go/bin:/data/service/code/bin:/data/service/go/bin:/root/bin"
        GOPROXY="https://goproxy.cn,direct"
    }
    parameters {
      gitParameter branchFilter: 'origin/(.*)', defaultValue: 'main', name: 'branch', type: 'PT_BRANCH', description: '请选择分支', sortMode: 'DESCENDING_SMART', quickFilterEnabled: true, listSize: "15"
    }

    stages{
        stage('build'){
          steps {
              sh 'rm -fr dist'
              sh "make build"
              sh 'chmod 755 sobotctl'
              sh 'mkdir dist/sobotctl/tmp -p; cp -a ./sobotctl ./config.yml ./tools  dist/sobotctl'
              sh 'cd dist/sobotctl/tools/redis; rm -f redis-shake; wget -O redis-shake http://artifactory.ops.sobot.tech/packages/redis/redis-shake-3.1.11'
              sh 'cd dist; tar -zcf sobotctl.tar.gz ./sobotctl'
              sh "coscmd upload dist/sobotctl.tar.gz deploy/${VERSION}/ --skipmd5"
              sh 'rm -fr dist'
              sh "echo ${COS_URl}/${VERSION}/sobotctl.tar.gz"
          }
        }
    }
}
