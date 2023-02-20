<template lang="pug">
  c-help-icon.ma-1(icon="help", right)
    v-layout(wrap)
      v-flex(xs12) {{ $t('login.connectionProtocols') }}
      ul
        li {{ $t('login.base') }}
        li {{ $t('login.LDAP') }}
</template>
