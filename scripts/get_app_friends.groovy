g.V().has('uid',XUID).outE('contact').has('has_app',true).as('ce').inV().as('inv').select('ce','inv').select('ce').valueMap().as('e').select('inv').valueMap().as('v').select('e','v')
