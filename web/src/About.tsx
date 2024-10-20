import './About.css';

function About() {
  return (
    <div className="about-container">
      <p>
        <b>Humor Arena</b> is based on the research paper{' '}
        <a href="https://arxiv.org/abs/2405.07280" target="_blank" rel="noopener noreferrer">
          "Humor Mechanics: Advancing Humor Generation with Multistep Reasoning"
        </a>{' '}
        by Alexey Tikhonov and Pavel Shtykovskiy.
      </p>

      <h2>Brief Summary</h2>
      <p>
        The authors explore the generation of one-liner jokes using multi-step reasoning.
        They developed a prototype for humor generation and conducted experiments to
        evaluate its effectiveness compared to human-created jokes and baseline models.
        Their findings show that multi-step reasoning consistently improves the quality
        of generated humor.
      </p>

      <h2>LaTeX Reference</h2>
      <pre className="latex-reference">
{`@article{tikhonov2024humor,
  title={Humor Mechanics: Advancing Humor Generation with Multistep Reasoning},
  author={Tikhonov, Alexey and Shtykovskiy, Pavel},
  journal={arXiv preprint arXiv:2405.07280},
  year={2024}
}`}
      </pre>

      <h2>Detailed Summary</h2>
      <p>
        In their paper, the authors address the challenge of generating high-quality,
        novel humor using large language models (LLMs). While LLMs have shown promise
        in natural language generation, producing truly creative and funny jokes remains
        difficult. The researchers propose a data-driven approach to reconstruct the
        mechanics of humor without relying on existing humor theories.
      </p>
      <p>
        Their method involves:
      </p>
      <ul>
        <li>
          Inferring a humor-generation policy directly from a dataset of jokes using
          LLMs in a zero-shot manner.
        </li>
        <li>
          Generating humor based on a seed topic and brainstorming associations to
          connect distant concepts.
        </li>
        <li>
          Evaluating the generated jokes through human annotation to assess funniness
          and novelty.
        </li>
      </ul>
      <p>
        The results indicate that this multi-step reasoning approach not only improves
        the funniness of the jokes but also increases their novelty compared to
        zero-shot GPT-4 outputs and human-generated jokes from Reddit datasets.
      </p>
      <p>
        The paper also discusses the subjectivity of humor and the importance of
        considering cultural preferences and ethical concerns in humor generation.
      </p>
      <p>
        For more details, you can read the full paper on{' '}
        <a href="https://arxiv.org/abs/2405.07280" target="_blank" rel="noopener noreferrer">
          arXiv
        </a>.
      </p>
    </div>
  );
}

export default About;