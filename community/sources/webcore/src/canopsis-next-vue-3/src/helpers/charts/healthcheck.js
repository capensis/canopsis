import { HEALTHCHECK_NETWORK_GRAPH_OPTIONS } from '@/constants'

export const getHealthcheckNodeRenderedPositionDiff = (factor = 1) =>
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSpace *
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor *
  factor

export const getHealthcheckInfoDiv = (message, hasError = false) =>
  `<div class="text-body-1 text-medium-emphasis pre-wrap ${hasError ? 'text-error' : ''}">${message}</div>`
