import copyText from './common/copyText'

export default function directive(app){
  app.directive('copyText', copyText)
}