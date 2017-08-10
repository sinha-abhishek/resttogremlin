g.V().has(UID_PROPERTY,XUID).outE(CONTACT_EDGE_LABEL).as('ce').inV().as('inv').select('ce','inv').select('ce').valueMap().as('e').select('inv').valueMap().as('v').select('e','v')
