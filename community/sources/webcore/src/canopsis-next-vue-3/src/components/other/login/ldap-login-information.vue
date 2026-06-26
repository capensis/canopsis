<template>
  <c-help-icon
    icon="help"
    color="secondary"
    right
  >
    <c-layout wrap>
      {{ $t('login.connectionProtocols') }}
      <ul>
        <li>{{ $t('login.base') }}</li>
        <li>{{ $t('login.LDAP') }}</li>
      </ul>
    </c-layout>
  </c-help-icon>
</template>
