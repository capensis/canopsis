<template lang="pug">
  c-help-icon(icon="help", right)
    v-layout(wrap)
      v-flex(xs12) {{ $t('login.connectionProtocols') }}
      ul
        li {{ $t('login.base') }}
        li {{ $t('login.LDAP') }}
</template>
