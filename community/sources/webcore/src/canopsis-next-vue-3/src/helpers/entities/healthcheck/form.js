import { HEALTHCHECK_ENGINES_NAMES } from '@/constants'

export const healthcheckParametersToForm = (healthcheckParameters = {}) => ({
  queue: healthcheckParameters.queue || {
    limit: 0,
    enabled: false,
  },
  messages: healthcheckParameters.messages || {
    limit: 0,
    enabled: false,
  },
  ...Object.values(HEALTHCHECK_ENGINES_NAMES).reduce((acc, engineName) => {
    const engineParameters = healthcheckParameters[engineName] || {}

    acc[engineName] = {
      enabled: !!engineParameters.enabled,
      minimal: engineParameters.minimal,
      optimal: engineParameters.optimal,
    }

    return acc
  }, {}),
})
