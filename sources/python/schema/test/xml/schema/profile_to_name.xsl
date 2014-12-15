<?xml version="1.0" encoding="ISO-8859-1"?>

<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                targetNamespace="profile>name:1.0">

<xsl:output method="xml" version="1.0" encoding="UTF-8" indent="yes"/>

<xsl:template match="/profile">
 <name><xsl:value-of select="name"/></name>
</xsl:template>

</xsl:stylesheet>


