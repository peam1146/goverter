package builder

import "github.com/dave/jennifer/jen"

type Map struct{}

func (*Map) Matches(source, target *Type) bool {
	return source.Map && target.Map
}

func (*Map) Build(gen Generator, ctx *MethodContext, sourceID *JenID, source, target *Type) ([]jen.Code, *JenID, *Error) {
	targetMap := ctx.Name(target.ID())
	key, value := ctx.Map()

	block, newKey, err := gen.Build(ctx, VariableID(jen.Id(key)), source.MapKey, target.MapKey)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapkey> " + source.MapKey.T.String(),
			TargetID:   "[]",
			TargetType: "<mapkey> " + target.MapKey.T.String(),
		})
	}
	valueStmt, valueKey, err := gen.Build(ctx, VariableID(jen.Id(value)), source.MapValue, target.MapValue)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "[]",
			SourceType: "<mapvalue> " + source.MapValue.T.String(),
			TargetID:   "[]",
			TargetType: "<mapvalue> " + target.MapValue.T.String(),
		})
	}
	block = append(block, valueStmt...)
	block = append(block, jen.Id(targetMap).Index(newKey.Code).Op("=").Add(valueKey.Code))

	stmt := []jen.Code{
		jen.Id(targetMap).Op(":=").Make(target.TypeAsJen(), jen.Len(sourceID.Code.Clone())),
		jen.For(jen.List(jen.Id(key), jen.Id(value)).Op(":=").Range().Add(sourceID.Code)).
			Block(block...),
	}

	return stmt, VariableID(jen.Id(targetMap)), nil
}
