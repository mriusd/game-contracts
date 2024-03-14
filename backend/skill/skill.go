// skill.go

package skill

import (

)

type Skill struct {
    SkillId             int     `json:"skillId"`
    Name                string  `json:"name"`
    ActiveDistance      int     `json:"activeDistance"`
    Multihit            bool    `json:"multihit"`
    AttackSuccessRate   int     `json:"attackSuccessRate"`
    HitAngle            int     `json:"hitAngle"`
    Disabled            bool    `json:"disabled"`
}


/*
    Damage colors:

    if isIgnoreDefence { light yelow }
    else if isExcellent { light green }
    else if isCritical { light blue }

    if double { display twice damage/2 }

*/

var Skills = map[int]Skill{
    0: {
        SkillId:           0,
        Name:              "Malee", 
        Multihit:          false,
        ActiveDistance:    1,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    1: {
        SkillId:           1,
        Name:              "Slash", 
        Multihit:          false,
        ActiveDistance:    1,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    2: {
        SkillId:           2,
        Name:              "Arrow", 
        Multihit:          false,
        ActiveDistance:    5,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    3: {
        SkillId:           3,
        Name:              "Tripple Shot", 
        Multihit:          true,
        ActiveDistance:    4,
        AttackSuccessRate: 100,
        HitAngle:          180,
    },
    4: {
        SkillId:           4,
        Name:              "Dark Spirits", 
        Multihit:          true,
        ActiveDistance:    20,
        AttackSuccessRate: 100,
        HitAngle:          360,
    },
}


func Get (k int) Skill {
    return Skills[k]
}