this_week_ts = Long.parseLong(XTS);
score_to_add=XSCORE.toInteger();
subject_to_add=XSUB;
uid = XUID;
score = 0;
println "score = "+score_to_add
cur_key = "key_"+this_week_ts+"_"+subject_to_add;
del_key = "key_" + (this_week_ts - 86400*1000*21)+"_"+subject_to_add;
println "cur="+cur_key+" del="+del_key
println "g.V().has('uid',$uid).property($cur_key, union(values($cur_key), constant($score_to_add)).sum()).properties($del_key).drop()"
g.V().has('uid',uid).property(cur_key, union(values(cur_key), constant(score_to_add)).sum()).properties(del_key).drop()
graph.tx().commit()
