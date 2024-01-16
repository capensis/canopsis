<template>
  <c-help-icon
    icon="help"
    right
  >
    <v-layout wrap>
      <v-flex xs12>
        {{ $t('login.connectionProtocols') }}
      </v-flex>
      <ul>
        <li>{{ $t('login.base') }}</li>
        <li>{{ $t('login.LDAP') }}</li>
      </ul>
    </v-layout>
  </c-help-icon>
</template>
