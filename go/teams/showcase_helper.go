package teams

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
)

func GetTeamShowcase(ctx context.Context, g *libkb.GlobalContext, teamname string) (ret keybase1.TeamShowcase, err error) {
	t, err := GetForTeamManagementByStringName(ctx, g, teamname, false)
	if err != nil {
		return ret, err
	}

	arg := apiArg(ctx, "team/get")
	arg.Args.Add("id", libkb.S{Val: t.ID.String()})

	var rt rawTeam
	if err := g.API.GetDecode(arg, &rt); err != nil {
		return ret, err
	}
	return rt.Showcase, nil
}

type memberShowcaseRes struct {
	Status      libkb.AppStatus `json:"status"`
	IsShowcased bool            `json:"is_showcased"`
}

func (c *memberShowcaseRes) GetAppStatus() *libkb.AppStatus {
	return &c.Status
}

func GetTeamAndMemberShowcase(ctx context.Context, g *libkb.GlobalContext, teamname string) (ret keybase1.TeamAndMemberShowcase, err error) {
	t, err := GetForTeamManagementByStringName(ctx, g, teamname, false)
	if err != nil {
		return ret, err
	}

	role, err := t.myRole(ctx)
	if err != nil {
		return ret, err
	}

	arg := apiArg(ctx, "team/get")
	arg.Args.Add("id", libkb.S{Val: t.ID.String()})

	var teamRet rawTeam
	if err := g.API.GetDecode(arg, &teamRet); err != nil {
		return ret, err
	}
	ret.TeamShowcase = teamRet.Showcase

	if role.IsOrAbove(keybase1.TeamRole_ADMIN) {
		arg = apiArg(ctx, "team/member_showcase")
		arg.Args.Add("tid", libkb.S{Val: t.ID.String()})

		var memberRet memberShowcaseRes
		if err := g.API.GetDecode(arg, &memberRet); err != nil {
			return ret, err
		}

		ret.IsMemberShowcased = memberRet.IsShowcased
	}

	return ret, nil
}

func SetTeamShowcase(ctx context.Context, g *libkb.GlobalContext, teamname string, isShowcased *bool, description *string) error {
	t, err := GetForTeamManagementByStringName(ctx, g, teamname, true)
	if err != nil {
		return err
	}

	if isShowcased == nil && description == nil {
		return errors.New("both isShowcased and description cannot be nil")
	}

	arg := apiArg(ctx, "team/team_showcase")
	arg.Args.Add("tid", libkb.S{Val: string(t.ID)})
	if isShowcased != nil {
		arg.Args.Add("is_showcased", libkb.B{Val: *isShowcased})
	}
	if description != nil {
		if len(*description) > 0 {
			arg.Args.Add("description", libkb.S{Val: *description})
		} else {
			arg.Args.Add("clear_description", libkb.B{Val: true})
		}
	}
	_, err = g.API.Post(arg)
	return err
}

func SetTeamMemberShowcase(ctx context.Context, g *libkb.GlobalContext, teamname string, isShowcased bool) error {
	t, err := GetForTeamManagementByStringName(ctx, g, teamname, true)
	if err != nil {
		return err
	}

	arg := apiArg(ctx, "team/member_showcase")
	arg.Args.Add("tid", libkb.S{Val: string(t.ID)})
	arg.Args.Add("is_showcased", libkb.B{Val: isShowcased})
	_, err = g.API.Post(arg)
	return err
}
