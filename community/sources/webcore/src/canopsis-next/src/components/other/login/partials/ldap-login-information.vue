<template>
  <c-help-icon
    icon="help"
    right="right"
  >
    <v-layout wrap="wrap">
      <v-flex xs12="xs12">
        {{ $t('login.connectionProtocols') }}
      </v-flex>
      <ul>
        <li>{{ $t('login.base') }}</li>
        <li>{{ $t('login.LDAP') }}</li>
      </ul>
    </v-layout>
  </c-help-icon>
</template>
