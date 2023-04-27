package repository

import (
	"bcraftTestTask/internal/models"
	"context"
	"database/sql"
	"strings"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db_ *sql.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db_,
	}
}

func (r *RecipeRepository) GetRecipes(ingredients, sortBy, time string) ([]models.Recipe, error) {
	var resultRecipes []models.Recipe
	var rows *sql.Rows
	var err error

	sort := ""

	if sortBy == "avgTime" {
		sort = "asc"
	} else if sortBy == "-avgTime" {
		sort = "desc"
	}

	if sort == "" {
		if time == "" {
			rows, err = r.db.Query("select step.recipeId, recipe.title, recipe.info, recipe.ingredients "+
				"from recipe left join step "+
				"on recipe.id = step.recipeId where ingredients like $1 group by step.recipeId, recipe.info, "+
				"recipe.title, recipe.ingredients;", "%"+strings.Replace(ingredients, ",", "%", -1)+"%")
		} else {
			rows, err = r.db.Query("select foo.recipeId, foo.title, foo.info, foo.ingredients from "+
				"(select step.recipeId, recipe.info, recipe.title, recipe.ingredients, sum(step.time) as time from "+
				"recipe left join step on recipe.id = step.recipeId group by step.recipeId, recipe.info, recipe.title,"+
				" recipe.ingredients) as foo where ingredients like $1 and foo.time = $2;",
				"%"+strings.Replace(ingredients, ",", "%", -1)+"%", time)
		}
	} else {
		if time == "" {
			rows, err = r.db.Query("select step.recipeId, recipe.title, recipe.info, recipe.ingredients "+
				"from recipe left join step "+
				"on recipe.id = step.recipeId where ingredients like $1 group by step.recipeId, recipe.title,"+
				" recipe.info, recipe.ingredients order by sum(step.time) "+sort+";",
				"%"+strings.Replace(ingredients, ",", "%", -1)+"%")
		} else {
			rows, err = r.db.Query("select foo.recipeId, foo.title, foo.info, foo.ingredients from "+
				"(select step.recipeId, recipe.info, recipe.title, recipe.ingredients, sum(step.time) as time "+
				" from recipe left join step on recipe.id = step.recipeId group by step.recipeId, recipe.info, "+
				"recipe.title, recipe.ingredients order by sum(step.time) "+sort+") as foo where ingredients "+
				"like $1 and foo.time = $2;", "%"+strings.Replace(ingredients, ",", "%", -1)+"%", time)
		}
	}

	defer closeRows(rows)
	if err != nil {
		return nil, err
	}

	var recipes []models.Recipe
	for rows.Next() {
		var stepRows *sql.Rows
		recipe := models.Recipe{}
		err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Info, &recipe.Ingredients)
		if err != nil {
			return nil, err
		}

		stepRows, err = r.db.Query("select step.id, step.info, step.time from step where recipeId = $1", recipe.Id)
		if err != nil {
			return nil, err
		}

		for stepRows.Next() {
			step := models.Step{}

			err := stepRows.Scan(&step.Id, &step.Info, &step.Time)
			if err != nil {
				return nil, err
			}
			recipe.Steps = append(recipe.Steps, step)
		}

		recipes = append(recipes, recipe)
		stepRows.Close()
	}

	for _, v := range recipes {
		resultRecipes = append(resultRecipes, v)
	}

	return resultRecipes, nil
}

func (r *RecipeRepository) GetRecipeById(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	rows, err := r.db.Query("select recipe.id, recipe.title, recipe.info, recipe.ingredients, "+
		"step.id, step.info, step.time from recipe left join step on recipe.id = step.recipeId where recipe.id = $1", id)
	defer closeRows(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		nullStep := models.DTOStep{}
		err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.Info, &recipe.Ingredients,
			&nullStep.Id, &nullStep.Info, &nullStep.Time)
		if err != nil {
			return nil, err
		}
		if !nullStep.Id.Valid {
			return &recipe, nil
		}

		recipe.Steps = append(recipe.Steps, models.Step{Id: nullStep.Id.Int64, Info: nullStep.Info.String,
			Time: nullStep.Time.Int64})
	}

	return &recipe, nil
}

func (r *RecipeRepository) PutRecipe(ctx context.Context, recipe models.Recipe) (int, error) {
	status := 0

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err != nil {
		return 0, err
	}

	if recipe.Id == 0 {
		var err error
		row := tx.QueryRow("insert into recipe (title, info, ingredients) "+
			"values ($1, $2, $3) RETURNING id;", recipe.Title, recipe.Info, recipe.Ingredients)

		if err != nil {
			return 0, err
		}
		err = row.Scan(&recipe.Id)
		if err != nil {
			return 0, err
		}

		status = 201
	} else {
		_, err := tx.Exec("update recipe set title = $1, info = $2, ingredients = $3 "+
			"where id = $4;", recipe.Title, recipe.Info, recipe.Ingredients, recipe.Id)

		if err != nil {
			return 0, err
		}

		status = 200
	}

	for _, step := range recipe.Steps {
		if step.Id == 0 {
			_, err := tx.Exec("insert into step (info, recipeId, time)"+
				" values ($1, $2, $3);", step.Info, recipe.Id, step.Time)

			if err != nil {
				return 0, err
			}
		} else {
			_, err := tx.Exec("update step set info = $1, time = $2"+
				" where id = $3 and recipeId = $4;", step.Info, step.Time, step.Id, recipe.Id)

			if err != nil {
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (r *RecipeRepository) DeleteRecipeById(id uint) (models.Recipe, error) {
	var recipe models.Recipe
	_, err := r.db.Exec("delete from recipe where recipe.id = $1", id)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}
