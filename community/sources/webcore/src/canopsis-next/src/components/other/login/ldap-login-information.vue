<template lang="pug">
  v-tooltip.ma-1(right)
    v-icon(slot="activator") help
    v-layout(wrap)
      v-flex(xs12) {{ $t('login.connectionProtocols') }}
      ul
        li {{ $t('login.base') }}
        li {{ $t('login.LDAP') }}
</template>
